import { mount } from "@cypress/react18";
import { Header } from "./Header";
import { createRouter, RouterContextProvider } from "@tanstack/react-router";
import { routeTree } from "@/routeTree.gen";

describe("Header", () => {
  it("should render the header contents", () => {
    const router = createRouter({ routeTree });
    mount(
      <>
        <RouterContextProvider router={router}>
          <Header />
        </RouterContextProvider>
      </>,
    );
    cy.contains("DevBlog").should("exist");
    cy.contains("Home").should("exist");
    cy.contains("About").should("exist");
  });
});
