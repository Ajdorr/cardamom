import { render, screen } from "@testing-library/react"
import AppCache from '@core/AppCache'
import { AppCacheContext } from '@core/AppCache'
import { api } from '@core/api'
import { useContext } from "react"
import { act } from "react-dom/test-utils"

function TestAppCacheComponent() {

    const { grocery, inventory } = useContext(AppCacheContext)

    return (<div>
        <div data-testid="grocery">
            {grocery.map(g => { return (<div key={g.uid} data-testid={g.uid}>{g.item}</div>) })}
        </div>
        <div data-testid="inventory">
            {inventory.map(i => { return (<div key={i.uid} data-testid={i.uid}>{i.item}</div>) })}
        </div>
    </div>)
}


it('App Cache', async () => {

    render(<AppCache><TestAppCacheComponent /></AppCache>);
    const [fulfilled, err] = (api.interceptors.response.use as jest.Mock).mock.calls[0]

    act(() => {
        // Reload cache
        fulfilled({ config: { url: "" } })

        fulfilled({
            config: { url: "grocery/create" },
            data: {
                uid: "abc",
                item: "potato",
            },
        })
    })

    expect((await screen.findByTestId("grocery")).children.length).toEqual(1)
    expect((await screen.findByTestId("abc")).innerHTML).toContain("potato")

    act(() => {
        fulfilled({
            config: { url: "grocery/create-batch" },
            data: [{
                uid: "abc",
                item: "potato",
            },
            {
                uid: "abc",
                item: "apple",
            }],
        })
    })
    expect((await screen.findByTestId("abc")).innerHTML).toContain("apple")

    act(() => {
        fulfilled({
            config: { url: "grocery/list" },
            data: [{
                uid: "grocery-0",
                item: "bread",
            },
            {
                uid: "grocery-1",
                item: "milk",
            },
            {
                uid: "grocery-2",
                item: "water",
            }],
        })
    })

    expect((await screen.findByTestId("grocery-0")).innerHTML).toContain("bread")
    expect((await screen.findByTestId("grocery-1")).innerHTML).toContain("milk")
    expect((await screen.findByTestId("grocery-2")).innerHTML).toContain("water")

    act(() => {
        fulfilled({
            config: {
                url: "grocery/delete",
                data: "{}",
            },
        })

        fulfilled({
            config: {
                url: "grocery/delete",
                data: JSON.stringify({
                    uid: "abc",
                    item: "potato",
                }),
            },
        })
    })

    expect(document.getElementById("abc")).toBeNull()

    act(() => {
        fulfilled({
            config: { url: "grocery/collect" },
            data: {
                grocery_item: {
                    uid: "abc",
                    item: "potato",
                },
                inventory_item: {
                    uid: "def",
                    item: "banana",
                    in_stock: true
                }
            },
        })
    })

    expect((await screen.findByTestId("grocery")).firstChild?.textContent).toContain("potato")
    expect((await screen.findByTestId("inventory")).firstChild?.textContent).toContain("banana")

    await act(() => {
        fulfilled({config: {url: ""}})
        fulfilled({
            config: { url: "inventory/create" },
            data: {
                uid: "xyz",
                item: "mushroom",
                in_stock: true
            },
        })
    })

    expect((await screen.findByTestId("inventory")).children.length).toEqual(1)
    expect((await screen.findByTestId("xyz")).innerHTML).toContain("mushroom")

    act(() => {
        fulfilled({
            config: { url: "inventory/create-batch" },
            data: [{
                uid: "xyz",
                item: "mushroom",
                in_stock: true,
            },
            {
                uid: "xyz",
                item: "grapes",
                in_stock: true,
            }],
        })
    })

    expect((await screen.findByTestId("xyz")).innerHTML).toContain("grapes")

    act(() => {
        fulfilled({
            config: { url: "grocery/clear" },
        })

        fulfilled({
            config: { url: "inventory/list" },
            data: [{
                uid: "inventory-0",
                item: "bread",
            },
            {
                uid: "inventory-1",
                item: "milk",
            },
            {
                uid: "inventory-2",
                item: "water",
            }],
        })
    })

    expect((await screen.findByTestId("inventory-0")).innerHTML).toContain("bread")
    expect((await screen.findByTestId("inventory-1")).innerHTML).toContain("milk")
    expect((await screen.findByTestId("inventory-2")).innerHTML).toContain("water")


    await act(() => { fulfilled({ config: { url: "grocery/clear" } }) })
    await act(() => { fulfilled({ config: { url: "inventory/delete" } }) })

    err("test").catch(() => {})
})