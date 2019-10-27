import akka.actor.ActorSystem
import akka.http.scaladsl.Http
import akka.http.scaladsl.model._
import akka.http.scaladsl.server.Directives._
import akka.http.scaladsl.server.directives.FileInfo
import akka.stream.{ActorMaterializer, Materializer}

import scala.concurrent.ExecutionContextExecutor
import scala.io.StdIn


object Server {

  def main(args: Array[String]): Unit = {
    implicit val system: ActorSystem = ActorSystem("media-server")
    implicit val materializer: ActorMaterializer = ActorMaterializer()
    implicit val executionContext: ExecutionContextExecutor = system.dispatcher

    val route = concat(
      path("upload") {
        concat(
          post {
            extractRequestContext { ctx =>
              implicit val materializer: Materializer = ctx.materializer
              fileUpload("fileUpload") {
                case (fileInfo: FileInfo, fileStream) =>
                  complete(HttpEntity(ContentTypes.`text/plain(UTF-8)`, fileInfo.fileName))
              }
            }
          }
        )
      },
      path("delete") {
        concat(
          delete {
            parameters(Symbol("url")) { (url) =>
              complete(HttpEntity(ContentTypes.`text/plain(UTF-8)`, url))
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
