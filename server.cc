#include "service.grpc.pb.h"
#include <cstdlib>
#include <grpc++/grpc++.h>
#include <grpcpp/server_context.h>
#include <grpcpp/support/status.h>
#include <iostream>
#include <opencv2/dnn.hpp>
#include <opencv2/imgproc.hpp>
#include <opencv2/objdetect.hpp>
#include <opencv2/objdetect/face.hpp>
#include <opencv2/opencv.hpp>
#include <string>

class FaceRecognitionServiceImpl : public FaceRecognitionService::Service {
private:
  grpc::HealthCheckServiceInterface *health_check_service_ = nullptr;
  const std::string model_path_ = "./models/yunet_n_320_320.onnx";
  cv::Ptr<cv::FaceDetectorYN> detector_ = cv::FaceDetectorYN::create(model_path_, "", cv::Size(320, 320), 0.8, 0.3, 5000);
  cv::Ptr<cv::FaceRecognizerSF> recognizer_ = cv::FaceRecognizerSF::create(model_path_, "");

  std::vector<std::vector<int>> find_combination(long n) {
    std::vector<std::vector<int>> result;
    for (int i = 0; i < n; i++) {
      for (int j = i + 1; j < n; j++) {
        result.push_back(std::vector<int>{i, j});
      }
    }
    return result;
  }

  cv::Mat get_face(const cv::Mat &faces) {
    int index = 0;
    for (int i = 1; i < faces.rows; i++) {
      auto w1 = faces.at<float>(index, 2);
      auto h1 = faces.at<float>(index, 3);
      auto w2 = faces.at<float>(i, 2);
      auto h2 = faces.at<float>(i, 3);

      if (w1 * h1 < w2 * h2)
        index = i;
    }
    return faces.row(index);
  }

  cv::Mat detect_faces(const std::string &data) {
    std::vector<uchar> raw_data(data.begin(), data.end());
    cv::Mat image = cv::imdecode(raw_data, cv::IMREAD_COLOR);
    this->detector_->setInputSize(image.size());
    cv::Mat faces;
    this->detector_->detect(image, faces);
    return faces;
  }

  grpc::Status RecognizeFaces(grpc::ServerContext *context, const RecognizeFacesRequest *request, RecognizeFacesResponse *response) override {
    std::cout << "\nReceived RecognizeFaces" << std::endl;

    cv::Mat faces = this->detect_faces(request->image());
    if (faces.empty()) {
      response->set_status(DetectFaceStatus_Enum_ENUM_NO_FACES);
      return grpc::Status::OK;
    }
    cv::Mat target_face = this->get_face(faces);
    std::cout << "Detected faces: " << faces.rows << " " << target_face << std::endl;

    for (const auto &face : request->faces()) {
      const auto &landmarks = face.landmarks();
      cv::Mat face_mat(1, landmarks.size(), CV_32F, (void *)landmarks.data());
      std::cout << "Face Mat: " << face_mat << std::endl;

      float score = this->recognizer_->match(face_mat, target_face);
      std::cout << "Matching face with score: " << score << std::endl;

      if (score < 0.5) {
        response->set_status(DetectFaceStatus_Enum_ENUM_NO_MATCH);
        return grpc::Status::OK;
      }
    }
    response->set_status(DetectFaceStatus_Enum_ENUM_OK);
    response->set_valid(true);

    return grpc::Status::OK;
  }

  grpc::Status DetectFaces(grpc::ServerContext *context, const DetectFacesRequest *request, DetectFacesResponse *response) override {
    std::cout << "\nReceived DetectFaces" << std::endl;

    std::vector<cv::Mat> detected_faces;
    for (const auto &data : request->images()) {
      cv::Mat faces = this->detect_faces(data);
      if (faces.empty()) {
        response->set_status(DetectFaceStatus_Enum_ENUM_NO_FACES);
        response->clear_faces();
        return grpc::Status::OK;
      }
      cv::Mat face = this->get_face(faces);
      std::cout << "Detected: " << face << std::endl;

      detected_faces.push_back(face);

      response->set_status(DetectFaceStatus_Enum_ENUM_OK);
      Face *face_result = response->add_faces();
      for (int j = 0; j < 15; j++) {
        face_result->add_landmarks(face.at<float>(0, j));
      }
    }

    auto n = detected_faces.size();
    auto combination = this->find_combination(n);
    for (const auto &pair : combination) {
      cv::Mat face1 = detected_faces.at(pair[0]);
      cv::Mat face2 = detected_faces.at(pair[1]);
      float score = this->recognizer_->match(face1, face2);

      std::cout << "Matching faces (" << pair[0] << "," << pair[1] << "): " << score << std::endl;

      if (score < 0.5) {
        response->set_status(DetectFaceStatus_Enum_ENUM_NO_MATCH);
        response->clear_faces();
        return grpc::Status::OK;
      }
    }

    response->set_status(DetectFaceStatus_Enum_ENUM_OK);
    return grpc::Status::OK;
  }

public:
  void set_health_check_service(grpc::HealthCheckServiceInterface *health_check_service) {
    health_check_service_ = health_check_service;
  }
};

int main() {
  const char *port_env = std::getenv("PORT");
  std::string port = port_env ? port_env : "8080";
  std::string addr = "0.0.0.0:" + port;
  std::string server_address(addr);
  FaceRecognitionServiceImpl service;

  grpc::EnableDefaultHealthCheckService(true);
  grpc::ServerBuilder builder;
  builder.SetMaxSendMessageSize(100 * 1024 * 1024);
  builder.SetMaxReceiveMessageSize(100 * 1024 * 1024);
  builder.SetMaxMessageSize(100 * 1024 * 1024);
  builder.AddListeningPort(server_address, grpc::InsecureServerCredentials());
  builder.RegisterService(&service);

  std::unique_ptr<grpc::Server> server(builder.BuildAndStart());
  service.set_health_check_service(server->GetHealthCheckService());
  std::cout << "Server listening on " << server_address << std::endl;
  server->Wait();
  return 0;
}
