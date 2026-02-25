import './index';
import { expect } from 'chai';
import { FavoritesDialogSk } from './favorites-dialog-sk';

import { setUpElementUnderTest } from '../../../infra-sk/modules/test_util';

describe('favorites-dialog-sk', () => {
  const newInstance = setUpElementUnderTest<FavoritesDialogSk>('favorites-dialog-sk');

  let element: FavoritesDialogSk;
  beforeEach(async () => {
    element = newInstance((_el: FavoritesDialogSk) => {
      // Place here any code that must run after the element is instantiated but
      // before it is attached to the DOM (e.g. property setter calls,
      // document-level event listeners, etc.).
    });
    await element.updateComplete;
  });

  it('renders for new', async () => {
    expect(element).to.not.be.null;
    element.open('12345', '', '', 'url1.com');
    await element.updateComplete;

    const nameInput = element.querySelector('input[placeholder="Name"]') as HTMLInputElement;
    const urlInput = element.querySelector('input[placeholder="URL"]') as HTMLInputElement;

    expect(nameInput).to.not.be.null;
    expect(nameInput.value).to.equal('');
    expect(urlInput.value).to.equal('url1.com');
  });

  it('renders for update', async () => {
    expect(element).to.not.be.null;
    element.open('', 'Fav', 'Fav Desc', 'url.com');

    // Wait for LitElement to update the DOM
    await element.updateComplete;

    const nameInput = element.querySelector('input[placeholder="Name"]') as HTMLInputElement;
    const descInput = element.querySelector('input[placeholder="Description"]') as HTMLInputElement;
    const urlInput = element.querySelector('input[placeholder="URL"]') as HTMLInputElement;

    expect(nameInput).to.not.be.null;
    expect(descInput).to.not.be.null;
    expect(urlInput).to.not.be.null;

    expect(nameInput.value).to.equal('Fav');
    expect(descInput.value).to.equal('Fav Desc');
    expect(urlInput.value).to.equal('url.com');
  });
});
