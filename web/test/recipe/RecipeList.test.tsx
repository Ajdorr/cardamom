
import { api } from "@core/api"
import App from "@core/App"
import { act, render, screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"

describe("Recipe List", () => {

    it("List", async () => {
        localStorage.setItem("csrf_token", "token")
        window.history.pushState({}, '', '/recipe/list')

            ;
        (api.post as jest.Mock).mockReturnValueOnce(Promise.resolve({
            data: [
                {
                    uid: "a",
                    is_trashed: false,
                    name: "Bolonese",
                    description: "A classic italian recipe",
                    meal: "Dinner",
                    instructions: "make it delicious",
                    ingredients: [
                        {
                            order: 0,
                            quantity: 1,
                            unit: "cups",
                            item: "penne",
                            optional: false,
                        },
                        {
                            order: 1,
                            quantity: 0.5,
                            unit: "cups",
                            item: "ground beef",
                            optional: false,
                        }
                    ]
                }
            ]
        }))
        render(<App />)

        expect(await screen.findByText("Bolonese")).toBeInTheDocument()

            ;
        (api.post as jest.Mock).mockReturnValueOnce(Promise.resolve({
            data: {
                uid: "a",
                is_trashed: true,
                name: "Bolonese",
                description: "A classic italian recipe",
                meal: "Dinner",
                instructions: "make it delicious",
                ingredients: [
                    {
                        order: 0,
                        quantity: 1,
                        unit: "cups",
                        item: "penne",
                        optional: false,
                    },
                    {
                        order: 1,
                        quantity: 0.5,
                        unit: "cups",
                        item: "ground beef",
                        optional: false,
                    }
                ]
            }
        }))
        await act(async () => { (await screen.findByAltText("trash recipe")).click() })

        expect((api.post as jest.Mock).mock.lastCall).toEqual(["recipe/update", { uid: "a", is_trashed: true }])
    })


})