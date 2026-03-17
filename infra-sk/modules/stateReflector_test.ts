// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import { assert } from 'chai';
import sinon from 'sinon';
import { stateReflector } from './stateReflector';

describe('stateReflector', () => {
  let pushStateSpy: sinon.SinonSpy;
  let replaceStateSpy: sinon.SinonSpy;

  beforeEach(() => {
    pushStateSpy = sinon.spy(window.history, 'pushState');
    replaceStateSpy = sinon.spy(window.history, 'replaceState');
  });

  afterEach(() => {
    pushStateSpy.restore();
    replaceStateSpy.restore();
  });

  it('uses pushState by default', async () => {
    let state = { foo: 'bar' };
    const stateHasChanged = stateReflector(
      () => state,
      (newState) => {
        state = newState as any;
      }
    );

    // Force loaded = true
    window.dispatchEvent(new PopStateEvent('popstate'));
    await Promise.resolve();

    state.foo = 'baz';
    stateHasChanged();

    assert.isTrue(pushStateSpy.calledOnce, 'pushState should have been called');
    assert.isTrue(replaceStateSpy.notCalled, 'replaceState should NOT have been called');
    assert.include(pushStateSpy.firstCall.args[2] as string, 'foo=baz');
  });

  it('uses replaceState when requested', async () => {
    let state = { alpha: 'beta' };
    const stateHasChanged = stateReflector(
      () => state,
      (newState) => {
        state = newState as any;
      },
      true /* replaceState */
    );

    // Force loaded = true
    window.dispatchEvent(new PopStateEvent('popstate'));
    await Promise.resolve();

    state.alpha = 'gamma';
    stateHasChanged();

    assert.isTrue(replaceStateSpy.calledOnce, 'replaceState should have been called');
    assert.isTrue(pushStateSpy.notCalled, 'pushState should NOT have been called');
    assert.include(replaceStateSpy.firstCall.args[2] as string, 'alpha=gamma');
  });
});
