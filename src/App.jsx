import * as React from 'react';
import { createBrowserRouter, RouterProvider } from "react-router";

import AppTheme from "./shared-themes/AppTheme";
import HomePage from "./pages/HomePage";
import SignInPage from "./pages/signin/SignInPage";
import SignUpPage from "./pages/signup/SignUpPage";
import DashboardPage from "./pages/dashboard/DashboardPage";

const themeComponents = {
}

const router = createBrowserRouter([
  { path: "/", element: <HomePage /> },
  { path: "/sign-in", element: <SignInPage /> },
  { path: "/sign-up", element: <SignUpPage /> },
  { path: "/dashboard", element: <DashboardPage /> },
])

export default function App(props) {
  return (
    <AppTheme {...props} themeComponents={themeComponents}>
      <RouterProvider router={router} />
    </AppTheme>
  );
}
