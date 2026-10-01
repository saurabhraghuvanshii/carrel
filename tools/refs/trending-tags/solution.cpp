#include <algorithm>
#include <unordered_map>
#include <vector>
using namespace std;

vector<int> trending(vector<vector<int>>& posts, int window, int k) {
    long long last = posts.back()[0];
    unordered_map<int, int> growth;
    for (auto& p : posts) {
        if (p[0] > last - window) growth[p[1]]++;
        else if (p[0] > last - 2LL * window) growth[p[1]]--;
    }
    vector<pair<int, int>> up;  // {-growth, tag}
    for (auto& [tag, g] : growth)
        if (g > 0) up.push_back({-g, tag});
    sort(up.begin(), up.end());
    vector<int> out;
    for (int i = 0; i < (int)up.size() && i < k; i++) out.push_back(up[i].second);
    return out;
}
