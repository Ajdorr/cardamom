
import { api } from "@core/api"
import App from "@core/App"
import { act, render, screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"

describe("Recipe Search", () => {

    it("Search", async () => {
        localStorage.setItem("csrf_token", "token")
        window.history.pushState({}, '', '/recipe/search')

        const dom = render(<App />)

            ;
        (api.post as jest.Mock).mockReturnValue(Promise.resolve({ data: [{ uid: "a" }] }))

        const nameInput = dom.container.querySelector(".recipe-search-menu-name input")
        await act(async () => { await userEvent.type(nameInput!, "pie{enter}") })
        expect((api.post as jest.Mock).mock.lastCall).toEqual(["recipe/search", {name: "pie"}])
        await act(async () => { 
            await userEvent.clear(nameInput!) 
            await userEvent.type(nameInput!, "{enter}") 
        })
        expect((api.post as jest.Mock).mock.lastCall).toEqual(["recipe/search", {name: "pie"}])

        const showAdvancedButton = dom.container.querySelector(".recipe-search-menu-show-advanced img")!
        await act(async () => { await userEvent.click(showAdvancedButton!) })

        await act(async () => { await userEvent.type(nameInput!, "cake{enter}") })
        expect((api.post as jest.Mock).mock.lastCall).toEqual(["recipe/search", {name: "cake", meal: null, description: null, ingredient: null}])

        const mealInput = dom.container.querySelector(".recipe-search-advanced-meal select")
        await act(async () => { await userEvent.selectOptions(mealInput!, "dinner") })
        expect((api.post as jest.Mock).mock.lastCall).toEqual(["recipe/search", {name: "cake", meal: "dinner", description: null, ingredient: null}])

        const descInput = dom.container.querySelector(".recipe-search-advanced-description input")
        await act(async () => { await userEvent.type(descInput!, "{enter}") })
        expect((api.post as jest.Mock).mock.lastCall).toEqual(["recipe/search", {name: "cake", meal: "dinner", description: null, ingredient: null}])

        await act(async () => { await userEvent.type(descInput!, "delicious{enter}") })
        expect((api.post as jest.Mock).mock.lastCall).toEqual(["recipe/search", {name: "cake", meal: "dinner", description: "delicious", ingredient: null}])

        const ingredientInput = dom.container.querySelector(".recipe-search-advanced-ingredient input")
        await act(async () => { await userEvent.type(ingredientInput!, "flour{enter}") })
        expect((api.post as jest.Mock).mock.lastCall).toEqual(["recipe/search", {name: "cake", meal: "dinner", description: "delicious", ingredient: "flour"}])

    })


})