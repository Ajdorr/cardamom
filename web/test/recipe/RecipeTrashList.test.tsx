
import { api } from "@core/api"
import App from "@core/App"
import { act, render, screen } from "@testing-library/react"

describe("Recipe Trash List", () => {

    it("Trash", async () => {
        localStorage.setItem("csrf_token", "token")
        window.history.pushState({}, '', '/recipe/trash')

            ;
        (api.post as jest.Mock).mockReturnValueOnce(Promise.resolve({ data: [{ uid: "a", name: "Lasagna", }] }))

        let dom
        await act(async () => { dom = render(<App />) })

            ;
        (api.post as jest.Mock).mockReturnValueOnce(Promise.resolve({ data: [] }))
        await act(async () => { (await screen.findByAltText("trash recipe")).click() }) 

        expect(await screen.findByText("Trash is empty")).toBeInTheDocument()

    })

})