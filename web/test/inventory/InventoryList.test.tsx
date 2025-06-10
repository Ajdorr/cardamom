import { api } from "@core/api"
import App from "@core/App"
import { act, render, screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"

describe("Inventory List", () => {

    it("List", async () => {
        localStorage.setItem("csrf_token", "token")
        window.history.pushState({}, '', '/inventory')
        const dom = render(<App />)

        const [fulfilled, err] = (api.interceptors.response.use as jest.Mock).mock.calls[1]
        expect(await screen.findByText("Nothing in your inventory")).toBeInTheDocument()

        await act(async () => {
            // Reload cache
            await fulfilled({ config: { url: "" } })

            await fulfilled({
                config: { url: "grocery/create" },
                data: {
                    uid: "1",
                    item: "potato",
                }
            })

            await fulfilled({
                config: { url: "inventory/list" },
                data: [
                    {
                        uid: "a",
                        item: "potato",
                        user_uid: "user",
                        quantity: "user",
                        store: "metro",
                        in_stock: false,
                        category: "cooking"
                    }, {
                        uid: "b",
                        item: "carrot",
                        user_uid: "user",
                        quantity: "user",
                        store: "costco",
                        in_stock: "spices",
                    }, {
                        uid: "c",
                        item: "onion",
                        user_uid: "user",
                        quantity: "user",
                        store: "costco",
                        in_stock: "sauces",
                    }
                ],
            })
        })

        const addInventoryInput = dom.container.querySelector(".inventory-list-add-item input")!
        await act(async () => {
            await userEvent.type(addInventoryInput, " {enter}")
            await userEvent.type(addInventoryInput, "onion{enter}")
        })

        expect((api.post as jest.Mock).mock.lastCall)
            .not.toContain("inventory/create")

        await act(async () => {
            await userEvent.type(addInventoryInput, "cheese{enter}")
        })
        expect((api.post as jest.Mock).mock.lastCall)
            .toContain("inventory/create")

        const inventoryItemInput = dom.container.querySelector(".inventory-item-input input")!
        await act(async () => {
            await userEvent.clear(inventoryItemInput)
            await userEvent.type(inventoryItemInput, "ginger{enter}")
        })

        // show more modal
        const showMoreBtn = (await screen.findAllByAltText("inventory-show-more"))[0]
        await act(async () => { showMoreBtn.click() })
        expect(dom.container.querySelector(".inventory-modal-root")).toBeInTheDocument()

        await act(async () => {
            const itemInput = await dom.container.querySelector(".inventory-modal-item input")
            await userEvent.clear(itemInput!)
            await userEvent.type(itemInput!, "garlic{enter}")
        })
        expect((api.post as jest.Mock).mock.lastCall)
            .toEqual(["inventory/update", { uid: "a", item: "garlic" }])

        await act(async () => {
            await userEvent.selectOptions(
                (await dom.container.querySelector(".inventory-modal-category select"))!,
                "spices")
        })
        expect((api.post as jest.Mock).mock.lastCall)
            .toEqual(["inventory/update", { uid: "a", category: "spices" }])

        await act(async () => {
            await userEvent.click(await screen.findByAltText("Add to grocery")) 
        })
        expect((api.post as jest.Mock).mock.lastCall)
            .toEqual(["grocery/create", { item: "potato" }])

        await act(async () => {
            await userEvent.click(await screen.findByAltText("Delete item")) 
        })
        expect((api.post as jest.Mock).mock.lastCall)
            .toEqual(["inventory/update", { uid: "a", in_stock: false }])

    })

    it("Cooking submenu", async() => {
        localStorage.setItem("csrf_token", "token")
        window.history.pushState({}, '', '/inventory/cooking')
        render(<App />)
    })

    it("Spices submenu", async() => {
        localStorage.setItem("csrf_token", "token")
        window.history.pushState({}, '', '/inventory/spices')
        render(<App />)
    })

    it("Sauces submenu", async() => {
        localStorage.setItem("csrf_token", "token")
        window.history.pushState({}, '', '/inventory/sauces')
        render(<App />)
    })

    it("non-perishables submenu", async() => {
        localStorage.setItem("csrf_token", "token")
        window.history.pushState({}, '', '/inventory/non-perishables')
        render(<App />)
    })
    
    it("Non-cooking submenu", async() => {
        localStorage.setItem("csrf_token", "token")
        window.history.pushState({}, '', '/inventory/non-cooking')
        render(<App />)
    })
    
    it("Non existant submenu redirect", async() => {
        localStorage.setItem("csrf_token", "token")
        window.history.pushState({}, '', '/inventory/invalid')
        render(<App />)
        expect(window.location.pathname).toEqual("/inventory")
    })
})