import {
  Database,
  PostgresConnector,
  DataTypes,
  Model,
} from "https://deno.land/x/denodb/mod.ts";
import { config } from "https://deno.land/x/dotenv/mod.ts";

import {
  User,
  Militant,
  Member,
  Document,
  Flight
} from "https://raw.githubusercontent.com/linuxing3/deno-game-monitor/develop/denodb/mock/CoreModels.ts"

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
export const postOptions = {
  host: "127.0.0.1",
  port: 5432,
  username: "spider",
  password: "20090909",
  database: "spider",
};
const connection = new PostgresConnector(postOptions);

const postdb = new Database(connection)
postdb.link([Movie, Article, User, Militant, Member, Document, Flight]);
await postdb.sync({ drop: false });
await postdb.close();
