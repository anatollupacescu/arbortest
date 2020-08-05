import { writable } from "svelte/store";

let cats = [
  {
    commit: "df423i",
    message: "some new feature",
    nodes: [
      { id: "group", group: 2, status: "pass" },
      { id: "test", group: 2, status: "pass" },
    ],
    links: [{ source: "test", target: "group", value: 3 }],
  },  {
    commit: "eabb33",
    message: "a fix for an old bug",
    nodes: [
      { id: "second", group: 2, status: "pass" },
      { id: "pair", group: 2, status: "pass" },
    ],
    links: [{ source: "second", target: "pair", value: 3 }],
  },
];

export const store = writable(cats);

export const current = writable({});
