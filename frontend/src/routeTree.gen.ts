/* prettier-ignore-start */

/* eslint-disable */

// @ts-nocheck

// noinspection JSUnusedGlobalSymbols

// This file is auto-generated by TanStack Router

import { createFileRoute } from "@tanstack/react-router";

// Import Routes

import { Route as rootRoute } from "./routes/__root";

// Create Virtual Routes

const AboutLazyImport = createFileRoute("/about")();
const IndexLazyImport = createFileRoute("/")();
const GroupsGroupIdLazyImport = createFileRoute("/groups/$groupId")();

// Create/Update Routes

const AboutLazyRoute = AboutLazyImport.update({
  path: "/about",
  getParentRoute: () => rootRoute,
} as any).lazy(() => import("./routes/about.lazy").then((d) => d.Route));

const IndexLazyRoute = IndexLazyImport.update({
  path: "/",
  getParentRoute: () => rootRoute,
} as any).lazy(() => import("./routes/index.lazy").then((d) => d.Route));

const GroupsGroupIdLazyRoute = GroupsGroupIdLazyImport.update({
  path: "/groups/$groupId",
  getParentRoute: () => rootRoute,
} as any).lazy(() =>
  import("./routes/groups/$groupId.lazy").then((d) => d.Route),
);

// Populate the FileRoutesByPath interface

declare module "@tanstack/react-router" {
  interface FileRoutesByPath {
    "/": {
      id: "/";
      path: "/";
      fullPath: "/";
      preLoaderRoute: typeof IndexLazyImport;
      parentRoute: typeof rootRoute;
    };
    "/about": {
      id: "/about";
      path: "/about";
      fullPath: "/about";
      preLoaderRoute: typeof AboutLazyImport;
      parentRoute: typeof rootRoute;
    };
    "/groups/$groupId": {
      id: "/groups/$groupId";
      path: "/groups/$groupId";
      fullPath: "/groups/$groupId";
      preLoaderRoute: typeof GroupsGroupIdLazyImport;
      parentRoute: typeof rootRoute;
    };
  }
}

// Create and export the route tree

export interface FileRoutesByFullPath {
  "/": typeof IndexLazyRoute;
  "/about": typeof AboutLazyRoute;
  "/groups/$groupId": typeof GroupsGroupIdLazyRoute;
}

export interface FileRoutesByTo {
  "/": typeof IndexLazyRoute;
  "/about": typeof AboutLazyRoute;
  "/groups/$groupId": typeof GroupsGroupIdLazyRoute;
}

export interface FileRoutesById {
  __root__: typeof rootRoute;
  "/": typeof IndexLazyRoute;
  "/about": typeof AboutLazyRoute;
  "/groups/$groupId": typeof GroupsGroupIdLazyRoute;
}

export interface FileRouteTypes {
  fileRoutesByFullPath: FileRoutesByFullPath;
  fullPaths: "/" | "/about" | "/groups/$groupId";
  fileRoutesByTo: FileRoutesByTo;
  to: "/" | "/about" | "/groups/$groupId";
  id: "__root__" | "/" | "/about" | "/groups/$groupId";
  fileRoutesById: FileRoutesById;
}

export interface RootRouteChildren {
  IndexLazyRoute: typeof IndexLazyRoute;
  AboutLazyRoute: typeof AboutLazyRoute;
  GroupsGroupIdLazyRoute: typeof GroupsGroupIdLazyRoute;
}

const rootRouteChildren: RootRouteChildren = {
  IndexLazyRoute: IndexLazyRoute,
  AboutLazyRoute: AboutLazyRoute,
  GroupsGroupIdLazyRoute: GroupsGroupIdLazyRoute,
};

export const routeTree = rootRoute
  ._addFileChildren(rootRouteChildren)
  ._addFileTypes<FileRouteTypes>();

/* prettier-ignore-end */

/* ROUTE_MANIFEST_START
{
  "routes": {
    "__root__": {
      "filePath": "__root.tsx",
      "children": [
        "/",
        "/about",
        "/groups/$groupId"
      ]
    },
    "/": {
      "filePath": "index.lazy.tsx"
    },
    "/about": {
      "filePath": "about.lazy.tsx"
    },
    "/groups/$groupId": {
      "filePath": "groups/$groupId.lazy.tsx"
    }
  }
}
ROUTE_MANIFEST_END */
