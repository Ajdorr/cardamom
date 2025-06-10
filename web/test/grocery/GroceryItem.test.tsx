import { describe, expect } from '@jest/globals';
import { fireEvent, render, screen } from '@testing-library/react';
import { api } from '@core/api'
import { act } from 'react-dom/test-utils';
import GroceryItem, { DragToDeleteTolerance } from '@pages/grocery/GroceryItem';
import userEvent from '@testing-library/user-event';


describe('Grocery item', () => {

  it('Input', async () => {
    const groceryItem = {
      uid: "a",
      created_at: "",
      updated_at: "",
      user_uid: "a",
      item: "potato",
      quantity: 1,
      store: "costco",
      is_collected: false,
    }
    const dom = render(<GroceryItem model={groceryItem} stores={["costco", "metro"]} />)

    const input = await dom.container.querySelector(".grocery-item-input input")
    await act(async () => {
      await userEvent.clear(input!)
      await userEvent.type(input!, "banana{enter}")
    })

    expect((api.post as jest.Mock).mock.lastCall)
      .toEqual(["grocery/update", { uid: "a", item: "banana" }])

    const dropdown = await dom.container.querySelector(".grocery-item-store select")
    await act(async () => {
      await userEvent.selectOptions(dropdown!, "metro")
    })

    expect((api.post as jest.Mock).mock.lastCall)
      .toEqual(["grocery/update", { uid: "a", store: "metro" }])

    const root = await dom.container.querySelector(".grocery-item-root")
    await act(async () => { await fireEvent.touchEnd(root!) })
    expect((api.post as jest.Mock).mock.lastCall)
      .not.toContain("grocery/delete")

    await act(async () => {
      await fireEvent.touchStart(root!, { touches: [{ clientX: -5}] })
    })

    await act(async () => {
      await fireEvent.touchMove(root!, { touches: [{ clientX: 10}] })
    })
    expect(await dom.container.querySelector(".grocery-item-collect-indicator")).toBeNull()
    expect(await dom.container.querySelector(".grocery-item-delete-indicator")).toBeNull()

    await act(async () => {
      await fireEvent.touchMove(root!, { touches: [{ clientX: -40}] })
    })
    expect(await dom.container.querySelector(".grocery-item-delete-indicator")).not.toBeNull()
    expect(await dom.container.querySelector(".grocery-item-collect-indicator")).toBeNull()

    await act(async () => {
      await fireEvent.touchMove(root!, { touches: [{ clientX: -50}] })
    })
    expect(await screen.findByAltText("delete indicator")).not.toBeNull()

    await act(async () => {
      await fireEvent.touchMove(root!, { touches: [{ clientX: DragToDeleteTolerance }] })
    })

    expect(await dom.container.querySelector(".grocery-item-delete-indicator")).toBeNull()
    expect(await dom.container.querySelector(".grocery-item-collect-indicator")).not.toBeNull()
    expect(root?.getAttribute("style")).toContain(`translateX(${DragToDeleteTolerance + 5}px)`)

    await act(async () => { await fireEvent.touchEnd(root!) })
    expect((api.post as jest.Mock).mock.lastCall)
      .toEqual(["grocery/delete", { uid: "a" }])
  });

})