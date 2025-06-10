import { describe } from '@jest/globals';
import { act, render, screen } from '@testing-library/react';
import { BrowserRouter, useNavigate } from 'react-router-dom';
import AuthIndex from '@auth/AuthIndex';
import { startOAuthLogin } from '@auth/authService';
import { api } from '@core/api';
import App from '@core/App';

jest.mock('@auth/authService', () => {
  const actual = jest.requireActual('@auth/authService')
  return {
    ...actual,
    startOAuthLogin: jest.fn(),
  }
})

describe('Auth Index', () => {

  it('Login - start', async () => {
    (api.post as jest.Mock).mockReturnValue(Promise.resolve({
      headers: { "x-csrf-token": "token" }
    }))
    window.history.pushState({}, '', '/auth/login')
    const dom = render(<App />)

    const loginBtn = await screen.findByText("Sign in with Google")
    expect(loginBtn).toBeInTheDocument()
    await (act(async () => { await loginBtn.click() }))

    expect((startOAuthLogin as jest.Mock)).toHaveBeenCalled()
  });

  it('Login - complete', async () => {
    (api.post as jest.Mock).mockReturnValue(Promise.resolve({
      headers: { "x-csrf-token": "token" }
    }))
    window.history.pushState({}, '', '/auth/oauth-return/google?code=a&state=b')
    const dom = render(<App />)

    expect(await screen.findByText("No grocery items in your list")).toBeInTheDocument()
  });

  it('Login - complete with forward', async () => {

    sessionStorage.setItem("auth_forward_pathname", "/inventory");
    (api.post as jest.Mock).mockReturnValue(Promise.resolve({
      headers: {
        "x-csrf-token": "token"
      }
    }))
    window.history.pushState({}, '', '/auth/oauth-return/google?code=a&state=b')
    render(<App />)
    expect(await screen.findByText("Nothing in your inventory")).toBeInTheDocument()
    expect(window.location.pathname).toEqual("/inventory")
  });

  it('Login failed', async () => {
    window.history.pushState({}, '', '/oauth-return/google')
    render(<BrowserRouter><AuthIndex /></BrowserRouter>)
    expect(await screen.findByText("Something went wrong, please try again later")).toBeInTheDocument()
  });

})