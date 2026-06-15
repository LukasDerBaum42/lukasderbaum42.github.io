+++
title = "Dice Dungeon"
description = "A dungeon crawler roguelike I started developing during autumn break 2025 as a small side project."
+++


# Dice Dungeon

Dice Dungeon is a dungeon crawler roguelike I started developing during autumn break 2025 as a small
side project. Originally, the plan was for it to be finished before the end of the break and to have
less than 1000 lines of code, but that plan didn't work out as planned.

The only thing that actually worked out is that it is my first game to be released.

When I started this project, I didn't really have any plans for what the game should be, so I went to
ChatGPT for inspiration, where it came up with the idea of a DnD-inspired text adventure game.

I liked the idea and added my own ideas to it. Those were to make it a "full" turn-based 2D dungeon
crawler with a lot of RPG mechanics.

After the original plan didn't work out, due to the project getting bigger than planned, I set the goal
to release a playable version before the end of 2025, which I did.

The reason it took longer than planned is that I didn't want to release it without some features, which
took longer than expected, partly because I'm lazy.

These features were:

- Items
- Randomly generated dungeons
- Traps
- Mod support
- A shop

Another reason is that I decided to try to make the codebase actually maintainable. This includes
storing important data as JSON files and refactoring the rendering part into its own file.

The data is stored as JSON files to make modding easy. All you have to do is create a JSON file that
follows the given structure and add it to the mod folder.

The refactoring was done to make the code cleaner and to make it easier for me to later add a proper GUI
for a potential Steam release.
Plans I have include: Neuro integration and better combat mechanics.

[If you're interested, the game is out on itch.io, just click here or look in the link section.](https://lukasderbaum42.itch.io/dice-dungeon)

[←Back](../)