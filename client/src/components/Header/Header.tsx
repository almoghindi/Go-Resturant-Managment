import React from "react";
import { Link } from "react-router-dom";
import "./Header.css";

export interface HeaderProps {
  title: string;
  logoUrl: string;
  isLoggedIn: boolean;
}

const Header: React.FC<HeaderProps> = ({ title, logoUrl, isLoggedIn }) => {
  return (
    <header className="header">
      <div className="header-content">
        <Link to="/" className="header-logo-title">
          <img src={logoUrl} alt="Logo" className="header-logo" />
          <h1 className="header-title">{title}</h1>
        </Link>
      </div>
      <nav className="header-nav">
        <ul>
          {isLoggedIn ? (
            <>
              <li>
                <Link to="/menu">Menu</Link>
              </li>
              <li>
                <Link to="/tables">Tables</Link>
              </li>
              <li>
                <Link to="/orders">Orders</Link>
              </li>
              <li>
                <Link to="/invoices">Invoices</Link>
              </li>
            </>
          ) : (
            <>
              <li>
                <Link to="/signup" className="signup-link">
                  Sign Up
                </Link>
              </li>
              <li>
                <Link to="/login" className="login-link">
                  Login
                </Link>
              </li>
            </>
          )}
        </ul>
      </nav>
    </header>
  );
};

export default Header;
