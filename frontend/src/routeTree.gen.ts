/* prettier-ignore-start */

/* eslint-disable */

// @ts-nocheck

// noinspection JSUnusedGlobalSymbols

// This file is auto-generated by TanStack Router

import { createFileRoute } from "@tanstack/react-router";

// Import Routes

import { Route as rootRoute } from "./routes/__root";
import { Route as BrowseImport } from "./routes/browse";

// Create Virtual Routes

const ProfileLazyImport = createFileRoute("/profile")();
const DiscordLazyImport = createFileRoute("/discord")();
const IndexLazyImport = createFileRoute("/")();
const GroupsGroupIdLazyImport = createFileRoute("/groups/$groupId")();

// Create/Update Routes

const ProfileLazyRoute = ProfileLazyImport.update({
  path: "/profile",
  getParentRoute: () => rootRoute,
} as any).lazy(() => import("./routes/profile.lazy").then((d) => d.Route));

const DiscordLazyRoute = DiscordLazyImport.update({
  path: "/discord",
  getParentRoute: () => rootRoute,
} as any).lazy(() => import("./routes/discord.lazy").then((d) => d.Route));

const BrowseRoute = BrowseImport.update({
  path: "/browse",
  getParentRoute: () => rootRoute,
} as any);

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
    "/browse": {
      id: "/browse";
      path: "/browse";
      fullPath: "/browse";
      preLoaderRoute: typeof BrowseImport;
      parentRoute: typeof rootRoute;
    };
    "/discord": {
      id: "/discord";
      path: "/discord";
      fullPath: "/discord";
      preLoaderRoute: typeof DiscordLazyImport;
      parentRoute: typeof rootRoute;
    };
    "/profile": {
      id: "/profile";
      path: "/profile";
      fullPath: "/profile";
      preLoaderRoute: typeof ProfileLazyImport;
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
  "/browse": typeof BrowseRoute;
  "/discord": typeof DiscordLazyRoute;
  "/profile": typeof ProfileLazyRoute;
  "/groups/$groupId": typeof GroupsGroupIdLazyRoute;
}

export interface FileRoutesByTo {
  "/": typeof IndexLazyRoute;
  "/browse": typeof BrowseRoute;
  "/discord": typeof DiscordLazyRoute;
  "/profile": typeof ProfileLazyRoute;
  "/groups/$groupId": typeof GroupsGroupIdLazyRoute;
}

export interface FileRoutesById {
  __root__: typeof rootRoute;
  "/": typeof IndexLazyRoute;
  "/browse": typeof BrowseRoute;
  "/discord": typeof DiscordLazyRoute;
  "/profile": typeof ProfileLazyRoute;
  "/groups/$groupId": typeof GroupsGroupIdLazyRoute;
}

export interface FileRouteTypes {
  fileRoutesByFullPath: FileRoutesByFullPath;
  fullPaths: "/" | "/browse" | "/discord" | "/profile" | "/groups/$groupId";
  fileRoutesByTo: FileRoutesByTo;
  to: "/" | "/browse" | "/discord" | "/profile" | "/groups/$groupId";
  id:
    | "__root__"
    | "/"
    | "/browse"
    | "/discord"
    | "/profile"
    | "/groups/$groupId";
  fileRoutesById: FileRoutesById;
}

export interface RootRouteChildren {
  IndexLazyRoute: typeof IndexLazyRoute;
  BrowseRoute: typeof BrowseRoute;
  DiscordLazyRoute: typeof DiscordLazyRoute;
  ProfileLazyRoute: typeof ProfileLazyRoute;
  GroupsGroupIdLazyRoute: typeof GroupsGroupIdLazyRoute;
}

const rootRouteChildren: RootRouteChildren = {
  IndexLazyRoute: IndexLazyRoute,
  BrowseRoute: BrowseRoute,
  DiscordLazyRoute: DiscordLazyRoute,
  ProfileLazyRoute: ProfileLazyRoute,
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
        "/browse",
        "/discord",
        "/profile",
        "/groups/$groupId"
      ]
    },
    "/": {
      "filePath": "index.lazy.tsx"
    },
    "/browse": {
      "filePath": "browse.tsx"
    },
    "/discord": {
      "filePath": "discord.lazy.tsx"
    },
    "/profile": {
      "filePath": "profile.lazy.tsx"
    },
    "/groups/$groupId": {
      "filePath": "groups/$groupId.lazy.tsx"
    }
  }
}
ROUTE_MANIFEST_END */
