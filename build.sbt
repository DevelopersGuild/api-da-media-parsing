name := "api-da-media-parsing"

version := "0.1"

scalaVersion := "2.13.1"

addSbtPlugin("au.com.onegeek" %% "sbt-dotenv" % "2.1.146")

libraryDependencies ++= Seq(
  "com.typesafe.akka" %% "akka-http"   % "10.1.10",
  "com.typesafe.akka" %% "akka-stream" % "2.5.23",
  "com.lightbend.akka" %% "akka-stream-alpakka-google-cloud-storage" % "1.1.2"
)