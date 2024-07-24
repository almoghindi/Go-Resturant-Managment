import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Header from "./components/Header/Header";
import { useAuth } from "./hooks/use-auth";
import logo from "./assets/logo.png";
import LandingPage from "./pages/LandingPage/LandingPage";
import SignupForm from "./components/Signup/SignupForm";

function App() {
  const { user } = useAuth();
  return (
    <Router>
      <Header title="Burger Palace" logoUrl={logo} isLoggedIn={false} />
      <Routes>
        <Route path="/" element={<LandingPage />} />
        {!user && <Route path="/signup" element={<SignupForm />} />}
      </Routes>
    </Router>
  );
}

export default App;
