#include <algorithm>
#include <unordered_map>
#include <vector>
using namespace std;

vector<int> topPages(vector<int>& views, int k) {
    unordered_map<int, int> count;
    for (int v : views) count[v]++;
    vector<pair<int, int>> pages;  // {-count, id}
    for (auto& [id, n] : count) pages.push_back({-n, id});
    sort(pages.begin(), pages.end());
    vector<int> out;
    for (int i = 0; i < (int)pages.size() && i < k; i++) out.push_back(pages[i].second);
    return out;
}
