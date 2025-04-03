import paho.mqtt.client as mqtt
import time

# MQTT broker settings
BROKER = "192.168.5.150"  # Your broker IP
PORT = 1883
TOPIC = "serverroom/door/#"  # Corrected topic with proper wildcard syntax

# Callback when client connects to broker
def on_connect(client, userdata, flags, rc):
    if rc == 0:
        print("✅ Connected to MQTT Broker!")
        try:
            # Subscribe to all door topics
            result, mid = client.subscribe(TOPIC)
            if result == mqtt.MQTT_ERR_SUCCESS:
                print(f"📡 Subscribed to topic: {TOPIC}")
            else:
                print(f"⚠ Failed to subscribe, return code: {result}")
        except ValueError as e:
            print(f"⚠ Subscription error: {e}")
    else:
        print(f"⚠ Failed to connect, return code {rc}")

# Callback when a message is received
def on_message(client, userdata, msg):
    try:
        door_number = msg.topic.split('/')[-1]  # Extract door number from topic
        status = msg.payload.decode('utf-8')    # Decode the message (open/closed)
        timestamp = time.strftime('%Y-%m-%d %H:%M:%S', time.localtime())
        
        # Print the status with timestamp
        emoji = "🚪" if status == "open" else "🔒"
        print(f"[{timestamp}] {emoji} Door {door_number}: {status.upper()}")
    except Exception as e:
        print(f"⚠ Error processing message: {e}")

# Callback when client disconnects
def on_disconnect(client, userdata, rc):
    if rc != 0:
        print("⚠ Unexpected disconnection. Attempting to reconnect...")

# Set up the MQTT client
client = mqtt.Client()
client.on_connect = on_connect
client.on_message = on_message
client.on_disconnect = on_disconnect

# Connect to the broker
try:
    print(f"🔌 Connecting to MQTT Broker at {BROKER}:{PORT}...")
    client.connect(BROKER, PORT, 60)
except Exception as e:
    print(f"⚠ Failed to connect to broker: {e}")
    exit(1)

# Start the loop to process network events
try:
    client.loop_forever()
except KeyboardInterrupt:
    print("\n👋 Disconnecting from broker...")
    client.disconnect()
    print("✅ Disconnected successfully")
except Exception as e: 
    print(f"⚠ Unexpected error: {e}")