import React from "react";
import { Link } from "react-router-dom";
import "./LandingPage.css";
import {
  FaUser,
  FaUtensils,
  FaTable,
  FaReceipt,
  FaClipboardList,
} from "react-icons/fa";

const LandingPage: React.FC = () => {
  return (
    <div className="landing-page">
      <section className="hero">
        <div className="hero-content">
          <h1>Welcome to Burger Palace</h1>
          <p>
            Manage your restaurant efficiently with our comprehensive toolset.
          </p>
          <div className="hero-buttons">
            <Link to="/menu" className="btn-primary">
              Manage your resturant now!
            </Link>
          </div>
        </div>
      </section>
      <section className="features">
        <h2>Features</h2>
        <div className="feature-list">
          <div className="feature-item">
            <FaUser size={32} />
            <h3>User Management</h3>
            <p>Manage user registrations, logins, and profiles easily.</p>
          </div>
          <div className="feature-item">
            <FaUtensils size={32} />
            <h3>Menu Management</h3>
            <p>Update and manage your restaurant's menu items effortlessly.</p>
          </div>
          <div className="feature-item">
            <FaTable size={32} />
            <h3>Table Management</h3>
            <p>Keep track of table availability and reservations.</p>
          </div>
          <div className="feature-item">
            <FaClipboardList size={32} />
            <h3>Order Management</h3>
            <p>Process and manage orders smoothly and efficiently.</p>
          </div>
          <div className="feature-item">
            <FaReceipt size={32} />
            <h3>Invoices</h3>
            <p>Generate and manage invoices with ease.</p>
          </div>
        </div>
      </section>
    </div>
  );
};

export default LandingPage;
