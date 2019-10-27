import akka.actor.ActorSystem
import akka.http.scaladsl.Http
import akka.http.scaladsl.model._
import akka.http.scaladsl.server.Directives._
import akka.http.scaladsl.server.directives.FileInfo
import akka.stream.alpakka.googlecloud.storage.scaladsl.GCStorage
import akka.stream.scaladsl.Source
import akka.stream.{ActorMaterializer, Materializer}
import akka.util.ByteString

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
                case (metadata: FileInfo, file: Source[ByteString, Any]) =>
                  println(file.toString())
                  val sink = GCStorage.resumableUpload("api-da-test-bucket", metadata.fileName, ContentTypes.`text/plain(UTF-8)`, )
                  complete(HttpEntity(ContentTypes.`text/plain(UTF-8)`, metadata.fileName))
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
