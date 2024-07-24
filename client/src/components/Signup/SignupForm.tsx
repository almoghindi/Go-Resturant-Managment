import React from "react";
import { useForm, SubmitHandler } from "react-hook-form";
import { useNavigate } from "react-router-dom";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import "./SignupForm.css";
import useRequest from "../../hooks/use-request";
import LoadingSpinner from "../loading-spinner";
import { useAuth } from "../../hooks/use-auth";

const schema = z.object({
  firstName: z.string().min(1, "First name is required"),
  lastName: z.string().min(1, "Last name is required"),
  phone: z.string().min(1, "Phone number is required"),
  email: z.string().email("Invalid email address"),
  password: z.string().min(6, "Password must be at least 6 characters long"),
});

type SignupFormData = z.infer<typeof schema>;

const SignupForm: React.FC = () => {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<SignupFormData>({
    resolver: zodResolver(schema),
  });
  const { sendRequest, isLoading, requestErrors } = useRequest();
  const { login } = useAuth();
  const navigate = useNavigate();

  const onSubmit: SubmitHandler<SignupFormData> = async (data) => {
    try {
      const response = await sendRequest({
        url: "http://localhost:8080/api/auth/signup",
        method: "POST",
        body: data,
        onSuccess: () => {
          navigate("/login");
        },
      });
      if (response) {
        const { userId, email, firstName, token, refreshToken } = response;
        login(userId, email, firstName, token, refreshToken);
        navigate("/dashboard");
      }
    } catch (error) {
      console.log(error);
    }
  };

  return (
    <>
      {isLoading && <LoadingSpinner />}
      <div className="signup-form-container">
        <h2>Sign Up</h2>
        <form onSubmit={handleSubmit(onSubmit)}>
          <div className="form-group">
            <label htmlFor="firstName">First Name</label>
            <input id="firstName" {...register("firstName")} />
            {errors.firstName && (
              <p className="error-message">{errors.firstName.message}</p>
            )}
          </div>
          <div className="form-group">
            <label htmlFor="lastName">Last Name</label>
            <input id="lastName" {...register("lastName")} />
            {errors.lastName && (
              <p className="error-message">{errors.lastName.message}</p>
            )}
          </div>
          <div className="form-group">
            <label htmlFor="phone">Phone</label>
            <input id="phone" {...register("phone")} />
            {errors.phone && (
              <p className="error-message">{errors.phone.message}</p>
            )}
          </div>
          <div className="form-group">
            <label htmlFor="email">Email</label>
            <input id="email" {...register("email")} />
            {errors.email && (
              <p className="error-message">{errors.email.message}</p>
            )}
          </div>
          <div className="form-group">
            <label htmlFor="password">Password</label>
            <input type="password" id="password" {...register("password")} />
            {errors.password && (
              <p className="error-message">{errors.password.message}</p>
            )}
          </div>
          <button type="submit" className="btn-submit">
            Sign Up
          </button>
          {requestErrors && (
            <p className="error-message">{requestErrors[0].message}</p>
          )}
        </form>
      </div>
    </>
  );
};

export default SignupForm;
