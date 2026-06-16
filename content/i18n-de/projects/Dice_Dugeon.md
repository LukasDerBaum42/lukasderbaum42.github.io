+++
title = "Dice Dungeon"
description = "Ein Dungeon-Crawler-Roguelike, dessen Entwicklung ich in den Herbstferien 2025 als kleines Nebenprojekt begonnen habe."
+++

# Dice Dungeon

Dice Dungeon ist ein Dungeon-Crawler-Roguelike, dessen Entwicklung ich in den Herbstferien 2025 als kleines Nebenprojekt begonnen habe. Ursprünglich war geplant, das Spiel noch vor dem Ende der Ferien fertigzustellen und dabei unter 1000 Zeilen Code zu bleiben, aber dieser Plan lief nicht ganz wie geplant.

Das Einzige, was tatsächlich funktioniert hat, ist, dass es mein erstes veröffentlichtes Spiel geworden ist.

Als ich dieses Projekt begonnen habe, hatte ich noch keine wirkliche Vorstellung davon, was das Spiel werden sollte. Deshalb habe ich ChatGPT um Inspiration gebeten, das die Idee eines von DnD inspirierten Text-Adventures vorgeschlagen hat.

Die Idee gefiel mir, und ich ergänzte sie um meine eigenen Vorstellungen. Daraus wurde schließlich ein „vollwertiger“ rundenbasierter 2D-Dungeon-Crawler mit vielen RPG-Mechaniken.

Nachdem der ursprüngliche Plan gescheitert war, weil das Projekt deutlich größer wurde als erwartet, setzte ich mir das Ziel, noch vor Ende 2025 eine spielbare Version zu veröffentlichen – was mir auch gelungen ist.

Der Grund, warum die Entwicklung länger dauerte als geplant, war, dass ich das Spiel nicht ohne bestimmte Funktionen veröffentlichen wollte. Deren Umsetzung dauerte länger als erwartet, teilweise auch, weil ich faul bin.

Diese Funktionen waren:

* Gegenstände
* Zufällig generierte Dungeons
* Fallen
* Mod-Unterstützung
* Ein Shop

Ein weiterer Grund war, dass ich versucht habe, die Codebasis tatsächlich wartbar zu gestalten. Dazu gehörte das Speichern wichtiger Daten als JSON-Dateien sowie das Auslagern des Renderings in eine eigene Datei.

Die Daten werden als JSON-Dateien gespeichert, um das Modding möglichst einfach zu machen. Alles, was man tun muss, ist eine JSON-Datei nach der vorgegebenen Struktur zu erstellen und sie in den Mod-Ordner zu legen.

Das Refactoring wurde durchgeführt, um den Code sauberer zu machen und es mir zu erleichtern, später für eine mögliche Steam-Veröffentlichung eine richtige grafische Benutzeroberfläche hinzuzufügen.

Zu meinen weiteren Plänen gehören eine Neuro-Integration sowie bessere Kampfmechaniken.

[Falls du interessiert bist: Das Spiel ist auf itch.io verfügbar. Klicke einfach hier oder schau im Link-Bereich nach.](https://lukasderbaum42.itch.io/dice-dungeon)

[← Zurück](../)
