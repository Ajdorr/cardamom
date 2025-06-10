import { describe, expect } from '@jest/globals';
import { render, screen } from '@testing-library/react';
import App from '@core/App'
import { api } from '@core/api'
import userEvent from '@testing-library/user-event';
import { act } from 'react-dom/test-utils';


describe('App', () => {

  it('smoke', () => {
    render(<App />)

    expect(document.querySelector(".home-root")).toBeInTheDocument()

  });

  it('login', () => {
    window.history.pushState({}, '', '/auth/login')
    render(<App />)

    expect(document.querySelector(".auth-index-root")).toBeInTheDocument()
  });

  it('redirect', async () => {
    window.history.pushState({}, '', '/grocery')
    render(<App />)
    expect(await screen.findByText("Sign in with Google")).toBeInTheDocument()
  });

  it('not found', async () => {
    localStorage.setItem("csrf_token", "token")
    window.history.pushState({}, '', '/zzz')
    var screen = render(<App />);
    (api.post as jest.Mock).mockReturnValue(Promise.resolve({}))

    expect(document.querySelector(".not-found-root")).toBeInTheDocument()
    expect(document.querySelector(".workspace-menu-bar-overlay")).toHaveStyle("display: none")

    var menuBtn = await screen.findByAltText("Show menu")
    await act(async () => menuBtn.click() )
    expect(document.querySelector(".workspace-menu-bar-overlay")).toHaveStyle("display: flex")

    await act(async () => userEvent.unhover(menuBtn))
    expect(document.querySelector(".workspace-menu-bar-overlay")).toHaveStyle("display: none")

    var groceryBtn = await screen.findByAltText("Go to grocery list")
    await act(async () => groceryBtn.click())
    expect(document.URL.endsWith("grocery"))

    var inventoryBtn = await screen.findByAltText("Go to inventory list")
    await act(async () => inventoryBtn.click())
    expect(document.URL.endsWith("inventory"))

    var recipeBtn = await screen.findByAltText("Go to recipe list")
    await act(async () => recipeBtn.click())
    expect(document.URL.endsWith("/recipe/list"))

    var accountBtn = await screen.findByAltText("Go to account settings")
    await act(async () => accountBtn.click())
    expect(document.URL.endsWith("account"))

  });

});