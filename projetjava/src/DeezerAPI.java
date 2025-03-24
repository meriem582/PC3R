import java.io.BufferedReader;
import java.io.InputStreamReader;
import java.net.HttpURLConnection;
import java.net.URL;
import org.json.JSONArray;
import org.json.JSONObject;
import java.net.URI;


public class DeezerAPI {
    public static void main(String[] args) {
        try {
            // URL de l'API Deezer
        URI uri = URI.create("https://api.deezer.com/search?q=gims");
URL url = uri.toURL();

            HttpURLConnection connection = (HttpURLConnection) url.openConnection();
            connection.setRequestMethod("GET");

            // Lire la réponse
            BufferedReader in = new BufferedReader(new InputStreamReader(connection.getInputStream()));
            String inputLine;
            StringBuilder response = new StringBuilder();
            while ((inputLine = in.readLine()) != null) {
                response.append(inputLine);
            }
            in.close();

            // Parser le JSON
            JSONObject jsonResponse = new JSONObject(response.toString());
            JSONArray tracks = jsonResponse.getJSONArray("data");

            // Afficher les informations sur chaque piste
            for (int i = 0; i < 5; i++) { // Limiter à 5 résultats
                JSONObject track = tracks.getJSONObject(i);
                System.out.println("Titre : " + track.getString("title"));
                System.out.println("Lien : " + track.getString("link"));
                System.out.println("Extrait : " + track.getString("preview"));
                System.out.println("Artiste : " + track.getJSONObject("artist").getString("name"));
                System.out.println("Album : " + track.getJSONObject("album").getString("title"));
                System.out.println("-----------------------------");
            }
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}