import akka.actor.ActorSystem
import akka.stream.{ActorMaterializer, Materializer}
import akka.http.scaladsl.Http
import akka.http.scaladsl.model._
import akka.http.scaladsl.server.Directives._
import akka.http.scaladsl.server.directives.FileInfo

import scala.concurrent.ExecutionContextExecutor
import scala.io.StdIn


object Server {
  def main(args: Array[String]): Unit = {
      implicit val system: ActorSystem = ActorSystem("media-server")
      implicit val materializer: ActorMaterializer = ActorMaterializer()
      implicit val executionContext: ExecutionContextExecutor = system.dispatcher

      val route = concat(
        path("upload"){
          concat(
            get {
              complete(HttpEntity(ContentTypes.`text/html(UTF-8)`, "<h1>Media Parsing Server</h1>"))
            },
            post {
              extractRequestContext{ ctx =>
                implicit val materializer: Materializer = ctx.materializer
                fileUpload("fileUpload") {
                  case(fileInfo: FileInfo, fileStream) =>
                    complete(HttpEntity(ContentTypes.`text/plain(UTF-8)`, fileInfo.fileName))
                }
              }
            }
          )
        }
    )

    val bindingFuture = Http().bindAndHandle(route, "localhost", 8080)
    println(s"Server online at http://localhost:8080/\nPress RETURN to stop...")
    StdIn.readLine() // let it run until user presses return
    bindingFuture
      .flatMap(_.unbind()) // trigger unbinding from the port
      .onComplete(_ => system.terminate()) // and shutdown when done

  }
}
