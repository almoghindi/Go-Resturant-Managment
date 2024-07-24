import { Meta, StoryFn } from "@storybook/react";
import { MemoryRouter } from "react-router-dom";
import LandingPage from "./LandingPage";

export default {
  title: "Components/LandingPage",
  component: LandingPage,
  decorators: [
    (Story) => (
      <MemoryRouter initialEntries={["/"]}>
        <Story />
      </MemoryRouter>
    ),
  ],
} as Meta;

const Template: StoryFn = (args) => <LandingPage {...args} />;

export const Default = Template.bind({});
Default.args = {};
