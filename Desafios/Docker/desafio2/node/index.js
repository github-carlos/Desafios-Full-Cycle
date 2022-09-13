const express = require("express");
const app = express();
const port = 3000;
const config = {
  host: "db",
  user: "root",
  password: "root",
  database: "nodedb",
};
const mysql = require("mysql");
const connection = mysql.createConnection(config);

const creatTable = `CREATE TABLE people(id int not null auto_increment, name varchar(255), primary key(id))`;
connection.query(creatTable, (_) => {
  const sql = `INSERT INTO people(name) values('Wesley')`;
  connection.query(sql);
  app.get("/", (req, res) => {
    const peopleSql = `SELECT * FROM people;`;
    connection.query(peopleSql, (err, result, fields) => {
      let listOfPeople = '';
      for (let i = 0; i < result.length; i++) {
        listOfPeople = listOfPeople.concat(`<li>${result[i]['name']}</li>`)
      }
      res.send(
        `
      </p> <p><h1>Full Cycle Rocks!</h1></p> <p>

      </p> <p>- Lista de nomes cadastrada no banco de dados.</p> <p> 
      <ul>
        ${listOfPeople}
      </ul>
      `
      );
    });
  });

  app.listen(port, () => {
    console.log("Rodando na porta " + port);
  });
});
