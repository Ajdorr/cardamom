import { describe, expect } from '@jest/globals';
import { render, screen } from '@testing-library/react';
import { api } from '@core/api'
import { act } from 'react-dom/test-utils';
import Account from '@auth/Account';
import { BrowserRouter } from 'react-router-dom';


describe('Account', () => {

  it('logout button', async () => {
    render(<BrowserRouter><Account /></BrowserRouter>)
    ;(api.post as jest.Mock).mockReturnValue(Promise.resolve({}))
    localStorage.setItem("csrf_token", "test_token")

    const btn = await screen.findByDisplayValue("Logout")
    await act(async () => { btn.click() })
    // expect((api.post as jest.Mock).mock.calls.length).toEqual(1)

    expect(localStorage.getItem("csrf_token")).toBeNull()
    expect(document.location.href).toBe('http://app.cardamom.cooking/auth/login')

  });

})