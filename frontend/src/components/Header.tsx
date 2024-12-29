import { Github } from "lucide-react";
import { ModeToggle } from "./mode-toggle";
import { Button } from "@/components/ui";
import { Link } from "@tanstack/react-router";

export const Header: React.FC = () => {
  return (
    <header className="flex items-center justify-between m-4">
      <h1 className="scroll-m-20 text-4xl font-extrabold tracking-tight lg:text-5xl">
        <span className="inline-block align-middle">
          <Link to="/">Rivals LFG</Link>
        </span>
      </h1>
      <div className="flex items-center space-x-2 p-2">
        <Link to="/" className="inline-flex items-center hover:underline">
          Home
        </Link>
        <Link
          to="/profile"
          className="inline-flex items-center hover:underline"
        >
          Profile
        </Link>
        <Button
          variant="ghost"
          size="icon"
          className="p-2"
          onClick={() => {
            window.open("https://github.com/jcserv/rivalslfg.com", "_blank");
          }}
        >
          <Github className="w-4 h-4" />
        </Button>
        <ModeToggle />
      </div>
    </header>
  );
};
