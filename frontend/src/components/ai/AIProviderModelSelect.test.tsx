import React from 'react';
import { act, create, type ReactTestRenderer } from 'react-test-renderer';
import { afterEach, describe, expect, it, vi } from 'vitest';

vi.mock('antd', () => ({
  Select: React.forwardRef((props: any, ref: any) => <><select ref={ref} {...props} />{props.open && props.dropdownRender?.(<div data-menu="options" />)}</>),
  Input: React.forwardRef((props: any, ref: any) => {
    const inputRef = React.useRef<HTMLInputElement>(null);
    // rc-input refreshes this handle on every render.
    React.useImperativeHandle(ref, () => ({ input: inputRef.current, focus: () => inputRef.current?.focus() }));
    return <input ref={inputRef} {...props} />;
  }),
}));
import AIProviderModelSelect from './AIProviderModelSelect';
import { t } from '../../i18n/catalog';

describe('searchable model selection', () => {
  let renderer: ReactTestRenderer | undefined;
  afterEach(async () => { await act(async () => renderer?.unmount()); vi.unstubAllGlobals(); });

  it('opens all candidates for a saved value, supports an explicit custom choice and clears search after selection', async () => {
    const changed = vi.fn();
    await act(async () => {
      renderer = create(<AIProviderModelSelect id="provider-model" aria-describedby="model-help" value="saved" onChange={changed} label="Model" placeholder="Choose" customLabel="Use custom:" options={[{ value: 'a', label: 'a' }, { value: 'b', label: 'b' }]} />);
    });
    const select = () => renderer!.root.findByType('select');
    expect(select().props.id).toBe('provider-model');
    expect(select().props['aria-describedby']).toBe('model-help');
    expect(select().props.searchValue).toBe('');
    expect(select().props.options.map((item: any) => item.value)).toEqual(['saved', 'a', 'b']);
    await act(async () => select().props.onSearch('custom-model'));
    expect(changed).not.toHaveBeenCalled();
    expect(select().props.options).toContainEqual({ value: 'custom-model', label: 'Use custom: custom-model' });
    await act(async () => select().props.onChange('custom-model'));
    expect(changed).toHaveBeenCalledWith('custom-model');
    expect(select().props.searchValue).toBe('');
    await act(async () => select().props.onChange(undefined));
    expect(changed).toHaveBeenLastCalledWith('');
  });

  it('does not lose the first typed character when keyboard input opens the dropdown', async () => {
    await act(async () => { renderer = create(<AIProviderModelSelect options={[]} label="Model" placeholder="Choose" customLabel="Use custom:" />); });
    const select = () => renderer!.root.findByType('select');
    await act(async () => { select().props.onSearch('x'); select().props.onOpenChange(true); });
    expect(select().props.searchValue).toBe('x');
    await act(async () => select().props.onBlur());
    expect(select().props.searchValue).toBe('');
  });

  it('manages enabled models while protecting default and SQL models, then updates normal choices', async () => {
    const changed = vi.fn();
    const Harness = () => {
      const [value, setValue] = React.useState('default');
      const [disabled, setDisabled] = React.useState(['hidden']);
      return <AIProviderModelSelect value={value} onChange={(next) => { changed(next); setValue(next); }}
        label="Model" placeholder="CLI default" customLabel="Use custom:" managementRequest={1}
        options={['default', 'fast', 'sql', 'hidden'].map((model) => ({ value: model, label: model }))}
        management={{ defaultModel: value, disabledModels: disabled, completionModel: 'sql', allowDefaultFallback: true,
          copy: (key, params) => t('en-US', key, params), source: 'Fixture catalog',
          onToggle: (model, enabled) => setDisabled((previous) => enabled ? previous.filter((item) => item !== model) : [...previous, model]), onAdd: vi.fn() }} />;
    };
    await act(async () => { renderer = create(<Harness />); });
    const select = () => renderer!.root.findByType('select');
    const toggle = (name: string) => renderer!.root.findByProps({ role: 'switch', 'aria-label': `Enable ${name}` });
    expect(toggle('default').props['aria-disabled']).toBe(true);
    expect(toggle('sql').props['aria-disabled']).toBe(true);
    await act(async () => { toggle('default').props.onClick(); toggle('sql').props.onClick(); });
    expect(toggle('default').props['aria-checked']).toBe(true);
    expect(toggle('sql').props['aria-checked']).toBe(true);
    await act(async () => toggle('fast').props.onClick());
    // Ant Design schedules this close after focus enters the custom popup.
    await act(async () => select().props.onOpenChange(false));
    expect(select().props.open).toBe(true);
    expect(toggle('fast').props['aria-checked']).toBe(false);
    expect(select().props.options.map((option: any) => option.value)).toEqual(['', 'default', 'sql']);
    expect(changed).not.toHaveBeenCalled();
    await act(async () => select().props.onChange('hidden'));
    expect(changed).not.toHaveBeenCalled();
    await act(async () => toggle('hidden').props.onClick());
    expect(select().props.options.map((option: any) => option.value)).toContain('hidden');
    await act(async () => renderer!.root.findByProps({ 'aria-label': 'Set default: hidden' }).props.onClick());
    expect(changed).toHaveBeenLastCalledWith('hidden');
    expect(toggle('hidden').props['aria-disabled']).toBe(true);
    expect(toggle('default').props['aria-disabled']).toBe(false);
    await act(async () => renderer!.root.findByProps({ 'aria-label': 'Close' }).props.onClick());
    expect(select().props.open).toBe(false);
  });

  it('adds custom candidates with Enter without replacing the default or duplicating a disabled name', async () => {
    const added = vi.fn(); const changed = vi.fn();
    await act(async () => { renderer = create(<AIProviderModelSelect value="default" onChange={changed}
      label="Model" placeholder="Choose" customLabel="Use custom:" managementRequest={1}
      options={['default', 'hidden'].map((model) => ({ value: model, label: model }))}
      management={{ defaultModel: 'default', disabledModels: ['hidden'], completionModel: '', allowDefaultFallback: false,
        copy: (key, params) => t('en-US', key, params), source: 'Saved catalog', onToggle: vi.fn(), onAdd: added }} />); });
    const search = () => renderer!.root.findByType('input');
    await act(async () => search().props.onChange({ target: { value: 'new-model' } }));
    await act(async () => search().props.onKeyDown({ key: 'Enter', stopPropagation: vi.fn(), preventDefault: vi.fn() }));
    expect(added).toHaveBeenCalledWith('new-model');
    expect(changed).not.toHaveBeenCalled();
    await act(async () => search().props.onChange({ target: { value: 'HIDDEN' } }));
    await act(async () => search().props.onKeyDown({ key: 'Enter', stopPropagation: vi.fn(), preventDefault: vi.fn() }));
    expect(added).toHaveBeenCalledTimes(1);
    await act(async () => search().props.onKeyDown({ key: 'Escape', stopPropagation: vi.fn() }));
    expect(renderer!.root.findByType('select').props.open).toBe(false);
  });

  it('focuses management once when visible without stealing switch focus on later renders', async () => {
    const focus = vi.fn();
    const node = { offsetWidth: 0, focus };
    let resize: () => void = () => undefined;
    const disconnect = vi.fn();
    vi.stubGlobal('ResizeObserver', class {
      constructor(callback: () => void) { resize = callback; }
      observe() {}
      disconnect = disconnect;
    });
    const renderSelect = (request: number) => <AIProviderModelSelect value="default" label="Model" placeholder="Choose" customLabel="Use custom:"
      managementRequest={request} options={[{ value: 'default', label: 'default' }]}
      management={{ defaultModel: 'default', disabledModels: [], completionModel: '', allowDefaultFallback: false,
        copy: (key, params) => t('en-US', key, params), source: 'Fixture catalog', onToggle: vi.fn(), onAdd: vi.fn() }} />;
    await act(async () => { renderer = create(renderSelect(1), { createNodeMock: (element) => element.type === 'input' ? node : null }); });
    expect(focus).not.toHaveBeenCalled();
    node.offsetWidth = 240;
    await act(async () => resize());
    expect(focus).toHaveBeenCalledTimes(1);
    expect(disconnect).toHaveBeenCalled();
    await act(async () => renderer!.root.findByProps({ role: 'switch' }).props.onClick());
    await act(async () => renderer!.root.findByType('input').props.onChange({ target: { value: 'd' } }));
    expect(focus).toHaveBeenCalledTimes(1);
    await act(async () => renderer!.root.findByProps({ 'aria-label': 'Close' }).props.onClick());
    await act(async () => renderer!.update(renderSelect(2)));
    expect(focus).toHaveBeenCalledTimes(2);
  });
});
