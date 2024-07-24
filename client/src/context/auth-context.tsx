import React, { createContext, ReactNode } from "react";
import { useAuth } from "../hooks/use-auth";

interface AuthContextType {
  userId: string;
  email: string;
  firstName: string;
  token: string;
  refreshToken: string;
  login: (
    userId: string,
    email: string,
    firstName: string,
    token: string,
    refreshToken: string
  ) => void;
  logout: () => void;
}

const defaultState: AuthContextType = {
  userId: "",
  email: "",
  firstName: "",
  token: "",
  refreshToken: "",
  login: () => {},
  logout: () => {},
};

export const AuthContext = createContext<AuthContextType>(defaultState);

export const AuthProvider: React.FC<{ children: ReactNode }> = ({
  children,
}) => {
  const { user, login, logout } = useAuth();

  const authContextValue = {
    userId: user?.user_id || "",
    email: user?.email || "",
    firstName: user?.firstName || "",
    token: user?.token || "",
    refreshToken: user?.refresh_token || "",
    login,
    logout,
  };

  return (
    <AuthContext.Provider value={authContextValue}>
      {children}
    </AuthContext.Provider>
  );
};
