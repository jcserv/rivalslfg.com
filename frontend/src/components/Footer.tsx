export const Footer: React.FC = () => {
  return (
    <footer>
      <div className="flex flex-col items-center justify-center m-8 text-center">
        <span>
          <p
            className="m-2 cursor-pointer"
            onClick={() => {
              window.open("https://jarrodservilla.com", "_blank");
            }}
          >
            Made with{" "}
            <span aria-label="heart" role="img">
              &#128153;
            </span>
            {" by Jarrod Servilla"}
          </p>
        </span>
        <span>
          <p className="text-sm text-muted-foreground">
            Rivals LFG is unofficial Fan Content, not approved/endorsed by Marvel or NetEase Games.
          </p>
          <p className="text-sm text-muted-foreground">
            Portions of the materials used are property of Marvel.
            © Marvel LLC.
          </p>
        </span>
      </div>
    </footer>
  );
};
