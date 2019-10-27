import java.io.InputStream

import akka.actor.ActorSystem
import akka.http.scaladsl.Http
import akka.http.scaladsl.model._
import akka.http.scaladsl.server.Directives._
import akka.http.scaladsl.server.directives.FileInfo
import akka.stream.scaladsl.{Source, StreamConverters}
import akka.stream.{ActorMaterializer, Materializer}
import akka.util.ByteString
import com.google.cloud.storage.{BlobId, BlobInfo, Storage, StorageOptions}

import scala.concurrent.ExecutionContextExecutor
import scala.io.{BufferedSource, StdIn, Source => ScalaSource}

object Server {

  def main(args: Array[String]): Unit = {
    implicit val system: ActorSystem = ActorSystem("media-server")
    implicit val materializer: ActorMaterializer = ActorMaterializer()
    implicit val executionContext: ExecutionContextExecutor = system.dispatcher

//    implicit val authFile: BufferedSource = ScalaSource.fromFile("../resources/auth.json")
//    authFile.close()
      sys.process.Process("env", None, "GOOGLE_APPLICATION_CREDENTIALS" -> "../resources/auth.json")

    val route = concat(
      path("upload") {
        concat(
          post {
            extractRequestContext { ctx =>
              implicit val materializer: Materializer = ctx.materializer
              fileUpload("fileUpload") {
                case (metadata: FileInfo, file: Source[ByteString, Any]) =>
//                  val storage: Storage = StorageOptions.getDefaultInstance.getService
//                  val bucketName: String = "api-da-test-bucket"
//                  val blobId = BlobId.of(bucketName, metadata.fileName)
//                  val blobInfo = BlobInfo.newBuilder(blobId).setContentType("text/plain").build
//                  val blob = storage.create(blobInfo, file) // Byte Array needed for file
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
