import {
  Database,
  PostgresConnector,
  DataTypes,
  Model,
} from "https://deno.land/x/denodb/mod.ts";
import { prompt, Input, Number, Confirm, Checkbox } from "https://deno.land/x/cliffy/prompt/mod.ts";

import {
  User,
  Militant,
  Member,
  Document,
  Flight
} from "https://raw.githubusercontent.com/linuxing3/deno-game-monitor/develop/denodb/mock/CoreModels.ts"
import { modelPool } from "https://raw.githubusercontent.com/linuxing3/deno-game-monitor/develop/denodb/mock/models.index.denodb.ts"

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

interface ModelPool {
  [key: string]: typeof Model
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

const extraTables: ModelPool = modelPool
const defaultTables: typeof Model[] = [Movie, Article, User, Militant, Member, Document, Flight]

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

async function menu() {

  const answers: PostOption = await prompt([{
    name: "tables",
    message: "Select some extra tables",
    type: Checkbox,
    options: Object.keys(extraTables),
  }, {
    name: "procede",
    message: "Do you really wan't to procede?",
    type: Confirm,
  }]);

  const mergedOptions = { ...answers, ...postOptions }
  console.table(mergedOptions)

  // TODO: choose tables model from list by table name
  if (answers.tables !== undefined) {
    answers.tables.forEach((v) => {
      console.log("Adding extra tables for you...")
      defaultTables.push(extraTables[v])
    })
  }
  if (answers.procede === true) {
    console.log("Creating tables for you...")
    createTable(mergedOptions, defaultTables)
  } else {
    console.log("Quit")
  }
}

await menu()
