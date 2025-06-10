
import { api } from "@core/api"
import App from "@core/App"
import RecipeSingle from "@pages/recipe/RecipeSingle"
import { act, render, screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"

describe("Recipe Shuffle", () => {

    beforeEach(() => {
        (RecipeSingle as jest.Mock).mockImplementation(() => { return (<div>Recipe Single</div>) })
    })

    it("Shuffle - positive flow", async () => {
        localStorage.setItem("csrf_token", "token")
        window.history.pushState({}, '', '/recipe/available')

            ;
        (api.post as jest.Mock).mockReturnValue(Promise.resolve({
            data: [
                {
                    uid: "a",
                    name: "Lasagna",
                    description: "Delicious!",
                    meal: "dinner",
                    ingredients: [
                        { uid: "a", quantity: "1", unit: "cups", item: "chicken" }
                    ],
                    instructions: "Combine and cook!"
                }
            ]
        }))

        let dom;
        await act(async () => { dom = await render(<App />) })
        const mealSelect = dom!.container.querySelector(".recipe-shuffle-form-meal select")
        await act(async () => { await userEvent.selectOptions(mealSelect, "None") })
        expect((api.post as jest.Mock).mock.lastCall).toEqual(["recipe/available", {meal: null}])

        await act(async () => { (await screen.findByDisplayValue("Yeah!")).click() })

        expect(window.location.pathname).toEqual("/recipe/edit/a")

    })

    it("Shuffle - negative flow", async () => {
        localStorage.setItem("csrf_token", "token")
        window.history.pushState({}, '', '/recipe/available')

            ;
        (api.post as jest.Mock).mockReturnValue(Promise.resolve({
            data: [
                {
                    uid: "a",
                    name: "Pesto Chicken",
                    description: "Delicious!",
                    meal: "dinner",
                    ingredients: [
                        { uid: "a", quantity: "1", unit: "cups", item: "chicken" }
                    ],
                    instructions: "Combine and cook!"
                },
                {
                    uid: "b",
                    name: "Lasagna",
                    description: "Delicious!",
                    meal: "dinner",
                    ingredients: [
                        { uid: "a", quantity: "1", unit: "cups", item: "chicken" }
                    ],
                    instructions: "Combine and cook!"
                }
            ]
        }))

        let dom;
        await act(async () => { dom = await render(<App />) })

        expect(await screen.findByText("Delicious!")).toBeInTheDocument()
        const recipeName = dom!.container.querySelector(".recipe-shuffle-name").textContent

        await act(async () => { (await screen.findByDisplayValue("Nah.")).click() })
        expect(dom!.container.querySelector(".recipe-shuffle-name").textContent).not.toEqual(recipeName)

        await act(async () => { (await screen.findByDisplayValue("Nah.")).click() })
        expect(dom!.container.querySelector(".recipe-shuffle-empty")).toBeInTheDocument()
        expect(await screen.findByText("That's all folks!")).toBeInTheDocument()

    })

})