import {
  Database,
  PostgresConnector,
  DataTypes,
  Model,
} from "https://deno.land/x/denodb/mod.ts";
import { config } from "https://deno.land/x/dotenv/mod.ts";
import { prompt, Input, Number, Confirm, Checkbox } from "https://deno.land/x/cliffy/prompt/mod.ts";

import {
  User,
  Militant,
  Member,
  Document,
  Flight
} from "https://raw.githubusercontent.com/linuxing3/deno-game-monitor/develop/denodb/mock/CoreModels.ts"
import { models } from "https://raw.githubusercontent.com/linuxing3/deno-game-monitor/develop/denodb/mock/models.index.denodb.ts"

interface PostOption {
  host: string;
  port: number;
  username: string;
  password: string;
  database: string;
  tables?: string[],
  procede?: boolean
}

class Movie extends Model {
  static table = "movies";
  static timestamps = true;

  static fields = {
    id: { primaryKey: true, autoIncrement: true },
    year: DataTypes.STRING,
    title: DataTypes.STRING,
    subtitle: DataTypes.STRING,
    desc: DataTypes.STRING,
    other: DataTypes.STRING,
    area: DataTypes.STRING,
    tag: DataTypes.STRING,
    star: DataTypes.STRING,
    comment: DataTypes.STRING,
    quote: DataTypes.STRING,
  };
}

class Article extends Model {
  static table = "articles";
  static timestamps = true;

  static fields = {
    id: { primaryKey: true, autoIncrement: true },
    title: DataTypes.STRING,
    url: DataTypes.STRING,
  };
}

const env = config({ safe: true });

const defaultTables = [Movie, Article, User, Militant, Member, Document, Flight]

const postOptions: PostOption = {
  host: "127.0.0.1",
  port: 5432,
  username: "spider",
  password: "20090909",
  database: "spider",
};

const createTable = async (options: any, tables: typeof Model[]) => {
  const postdb = new Database(new PostgresConnector(options))
  postdb.link(tables);
  await postdb.sync({ drop: false });
  await postdb.close();
}

const deleteTable = async (options: any, tables: typeof Model[]) => {
  const postdb = new Database(new PostgresConnector(options))
  postdb.link(tables);
  await postdb.sync({ drop: false });
  await postdb.close();
}

async function fullMenu() {

  const answers: PostOption = await prompt([{
    name: "host",
    message: "What's your postgresql host?",
    type: Input,
  }, {
    name: "port",
    message: "What's the port?",
    type: Number,
  }, {
    name: "username",
    message: "What's the username?",
    type: Input,
  }, {
    name: "password",
    message: "What's the password?",
    type: Input,
  }, {
    name: "database",
    message: "what's the database name?",
    type: Input,
  }, {
    name: "tables",
    message: "Select some tables",
    type: Checkbox,
    options: ["documents", "users", "articles", "movies"],
  }, {
    name: "procede",
    message: "Do you really wan't to procede?",
    type: Confirm,
  }]);

  const mergedOptions = { ...answers, ...postOptions }
  console.table(mergedOptions)

  if (answers.procede === true) {
    console.log("Creating tables for  you...")
    createTable(mergedOptions, defaultTables)
  } else {
    console.log("Deleting tables for  you...")
    deleteTable(mergedOptions, defaultTables)
  }
}

async function menu() {

  const answers: PostOption = await prompt([{
    name: "tables",
    message: "Select some tables",
    type: Checkbox,
    options: ["documents", "users", "articles", "movies"],
  }, {
    name: "procede",
    message: "Do you really wan't to procede?",
    type: Confirm,
  }]);

  const mergedOptions = { ...answers, ...postOptions }
  console.table(mergedOptions)

  // TODO: choose tables model from list by table name
  if (answers.procede === true) {
    console.log("Creating tables for  you...")
    createTable(mergedOptions, defaultTables)
  } else {
    console.log("Deleting tables for  you...")
    deleteTable(mergedOptions, defaultTables)
  }
}

await menu()
