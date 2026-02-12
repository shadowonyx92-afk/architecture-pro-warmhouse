require 'sinatra'
require 'json'
require 'thread'

set :bind, '0.0.0.0'
set :port, 8085

# Структура устройства
device = {
  id: "1",
  name: "Living Room Light",
  status: "off"
}

mutex = Mutex.new

# GET состояние
get '/light/:id' do
  if params[:id] != device[:id]
    status 404
    return { error: "Device not found" }.to_json
  end

  mutex.synchronize do
    content_type :json
    device.to_json
  end
end

# POST изменить статус
post '/light/:id/set' do
  if params[:id] != device[:id]
    status 404
    return { error: "Device not found" }.to_json
  end

  begin
    payload = JSON.parse(request.body.read, symbolize_names: true)
  rescue JSON::ParserError => e
    status 400
    return { error: e.message }.to_json
  end

  unless ["on", "off"].include?(payload[:status])
    status 400
    return { error: "Invalid status, must be 'on' or 'off'" }.to_json
  end

  mutex.synchronize do
    device[:status] = payload[:status]
    content_type :json
    device.to_json
  end
end
