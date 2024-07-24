console.log("Main")

import van from "van"

const { a, div, p, ul, li } = van.tags

const Hello = () => div(
    p("👋Hello"),
    ul(
        li("🗺️World"),
        li(a({href: "https://vanjs.org/"}, "🍦VanJS"))
    ),
)
  
van.add(document.body, Hello())