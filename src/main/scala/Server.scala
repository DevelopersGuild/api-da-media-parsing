import akka.actor.ActorSystem
import akka.stream.ActorMaterializer

import scala.concurrent.ExecutionContextExecutor

object Server {
  def main(args: Array[String]): Unit = {
      implicit val system: ActorSystem = ActorSystem("media-server")
      implicit val materializer: ActorMaterializer = ActorMaterializer()
      implicit val executionContext: ExecutionContextExecutor = system.dispatcher
  }
}
