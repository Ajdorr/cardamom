import { api } from '@core/api'
import AppCache from "@core/AppCache"
import GroceryList from "@pages/grocery/GroceryList"
import { act, render, screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"

describe("Grocery List", () => {

    it("List", async () => {
        const dom = render(<AppCache><GroceryList /></AppCache>)
        const [fulfilled, err] = (api.interceptors.response.use as jest.Mock).mock.calls[0]

        expect(await screen.findByText("No grocery items in your list")).toBeInTheDocument()
        await act(async () => {
            // Reload cache
            fulfilled({ config: { url: "" } })

            fulfilled({
                config: { url: "grocery/list" },
                data: [
                    {
                        uid: "a",
                        item: "potato",
                        user_uid: "user",
                        quantity: "user",
                        store: "metro",
                        is_collected: false,
                    }, {
                        uid: "b",
                        item: "carrot",
                        user_uid: "user",
                        quantity: "user",
                        store: "costco",
                        is_collected: false,
                    }, {
                        uid: "c",
                        item: "onion",
                        user_uid: "user",
                        quantity: "user",
                        store: "costco",
                        is_collected: true,
                    }
                ],
            })
        })

        expect(await screen.findByDisplayValue("potato")).toBeInTheDocument()

        const clearBtn = await screen.findByAltText("clear")
        await act(async () => await clearBtn.click())
        expect((api.post as jest.Mock).mock.lastCall).toContain("grocery/clear")

        const collectBtns = await screen.findAllByAltText("collect")
        await act(async () => await collectBtns[0].click())
        expect((api.post as jest.Mock).mock.lastCall)
            .toEqual(["grocery/collect", { uid: "a", is_collected: true }])

        const undoBtn = await screen.findByAltText("undo")
        await act(async () => await undoBtn.click())
        expect((api.post as jest.Mock).mock.lastCall)
            .toEqual(["grocery/collect", { uid: "c", is_collected: false }])

        const storeInputText = await screen.findByPlaceholderText("Add or select store")
        await act(async () => {
            await userEvent.type(storeInputText, "costco{enter}")
        })
        const grocery = await dom.container.querySelectorAll(".grocery-item-root")
        expect(grocery).toHaveLength(1)

        const itemInput = grocery[0].getElementsByClassName("grocery-item-input")[0]
        await act(async () => {
            await userEvent.type(itemInput, "costco{enter}")
        })

        const addGrocery = await screen.findByPlaceholderText("Add a grocery")
        await act(async () => {
            await userEvent.type(addGrocery, " {enter}")
            await userEvent.type(addGrocery, "pepper{enter}")
        })
        expect((api.post as jest.Mock).mock.lastCall).toContain("grocery/create")

    })
})