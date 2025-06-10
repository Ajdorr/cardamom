// jest-dom adds custom jest matchers for asserting on DOM nodes.
// allows you to do things like:
// expect(element).toHaveTextContent(/react/i)
// learn more: https://github.com/testing-library/jest-dom
import '@testing-library/jest-dom';


jest.mock("@core/api", () => ({
  __esModule: true,
  api: {
    get: jest.fn(),
    post: jest.fn(),
    interceptors: {
      request: { use: jest.fn(), eject: jest.fn() },
      response: { use: jest.fn(), eject: jest.fn() },
    }
  }
}
));


const mock_RecipeSingle = jest.fn()
jest.mock("@pages/recipe/RecipeSingle", () => ({
  __esModule: true,
  default: mock_RecipeSingle,
}
))

beforeEach(() => {
  mock_RecipeSingle.mockImplementation(jest.requireActual("@pages/recipe/RecipeSingle").default)
})
