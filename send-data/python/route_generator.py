import requests
import polyline
from datetime import datetime, timedelta
import json
import random

def add_jitter(coords, jitter_range=0.0001):
    jittered = []
    for lat, lon in coords:
        new_lat = lat + random.uniform(-jitter_range, jitter_range)
        new_lon = lon + random.uniform(-jitter_range, jitter_range)
        jittered.append((new_lat, new_lon))
    return jittered


def get_route_points(origin, destination, api_key, num_points=20):
    # Construct the URL for the Directions API request using lat, long instead of place names
    url = f"https://maps.googleapis.com/maps/api/directions/json?origin={origin}&destination={destination}&key={api_key}"

    # Send the GET request
    response = requests.get(url)
    
    if response.status_code == 200:
        directions = response.json()
        print(json.dumps(directions, indent=4))
        # Initialize a list to hold all the points
        all_points = []

        # Loop through each route in the directions
        for route in directions['routes']:
            for leg in route['legs']:
                # Loop through each step in the leg
                for step in leg['steps']:
                    # Get the polyline points for each step
                    encoded_polyline = step['polyline']['points']
                    decoded_points = polyline.decode(encoded_polyline)
                    all_points.extend(decoded_points)

        # Down-sample the points to get only `num_points` evenly spaced points
        if len(all_points) <= num_points:
            return all_points  # If there are fewer points than requested, return all

        # Calculate the step size to pick points evenly across the route
        step_size = len(all_points) // num_points

        # Select every `step_size`-th point
        downsampled_points = [all_points[i] for i in range(0, len(all_points), step_size)]

        # If the number of downsampled points is less than `num_points`, add the last point
        if len(downsampled_points) < num_points:
            downsampled_points.append(all_points[-1])

        return downsampled_points
    else:
        print("Error fetching data:", response.status_code)
        return None





if __name__ == "__main__":
    origin = "-12.123837, -76.949747"       
    destination = "-12.239833, -76.927514" 
    
    
    points = get_route_points(destination, origin, api_key, 300)
    with open("coordinates.txt", "w") as f:
        for lat, lon in points:
            f.write(f"{lat}, {lon}\n")

# coords_from_file = []
# with open("coordinates.txt", "r") as f:
#     for line in f:
#         lat_str, lon_str = line.strip().split(",")
#         coords_from_file.append((float(lat_str), float(lon_str)))
# print(coords_from_file)

