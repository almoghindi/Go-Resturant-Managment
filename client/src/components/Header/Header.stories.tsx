import { StoryFn, Meta } from "@storybook/react";
import { MemoryRouter } from "react-router-dom";
import Header, { HeaderProps } from "./Header";
import logo from "../../assets/logo.png";

export default {
  title: "Components/Header",
  component: Header,
  decorators: [
    (Story) => (
      <MemoryRouter>
        <Story />
      </MemoryRouter>
    ),
  ],
} as Meta;

const Template: StoryFn<HeaderProps> = (args) => <Header {...args} />;

export const Default = Template.bind({});
Default.args = {
  title: "Burger Palace",
  logoUrl: logo,
};
