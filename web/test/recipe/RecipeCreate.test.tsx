
import { api } from "@core/api"
import App from "@core/App"
import RecipeSingle from "@pages/recipe/RecipeSingle"
import { act, render, screen } from "@testing-library/react"

describe("Recipe Create", () => {

    beforeEach(() => {
        (RecipeSingle as jest.Mock).mockImplementation(() => { return (<div>Recipe Single</div>)})
    })

    it("Create", async () => {
        localStorage.setItem("csrf_token", "token")
        window.history.pushState({}, '', '/recipe/create')

        render(<App />)

            ;
        (api.post as jest.Mock).mockReturnValue(Promise.resolve({ data: { uid: "a", ingredients: [], instructions: "abc" } }))
        await act(async () => { (await screen.findByText("Create!")).click() })

        expect((api.post as jest.Mock).mock.calls[2]).toEqual(["recipe/create", { name: "" }])
        expect(window.location.pathname).toEqual("/recipe/edit/a")
    })

})