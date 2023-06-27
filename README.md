# Heimdall

A simple auth server for the rest of us.

## Why does this exist?

Whenever I get the idea for a new side project that would be hosted online, I'm immediately hit with the realization that I need auth in order to secure it.
Of course, I could skip the auth and just hope that my little corner of the internet never gets targeted, but security by obscurity is rarely a good idea.
I could reimplement a simple auth mechanism in each app, but updating this code everyone would eventually get annoying.
I'm a lazy programmer so a reusable auth server is the obvious choice.
But then comes the challenge of choosing, deploying, and configuring one of these beasts.
Popular auth servers like Keycloak are powerful, but challenging to work with and complete overkill for my little projects.

This is where Heimdall comes in.
Heimdall fills the void between no auth, and professional level auth servers.
It provides the tools needed to protect an application without requiring layers of configuration.
Whether it's users logging in or clients making server-server requests, Heimdall makes it easy to protect your hobby projects from the rest of the internet.

## Features

- User registration and login
- JWT based authentication
- Token introspection
- API Keys
