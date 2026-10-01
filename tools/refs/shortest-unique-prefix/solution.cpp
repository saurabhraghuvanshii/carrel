#include <array>
#include <string>
#include <vector>
using namespace std;

vector<string> shortestPrefixes(vector<string>& words) {
    vector<array<int, 26>> next(1);
    vector<int> pass(1, 0);
    next[0].fill(-1);
    for (auto& w : words) {
        int cur = 0;
        for (char c : w) {
            if (next[cur][c - 'a'] < 0) {
                next[cur][c - 'a'] = next.size();
                next.emplace_back();
                next.back().fill(-1);
                pass.push_back(0);
            }
            cur = next[cur][c - 'a'];
            pass[cur]++;
        }
    }
    vector<string> out;
    for (auto& w : words) {
        int cur = 0;
        for (size_t k = 0; k < w.size(); k++) {
            cur = next[cur][w[k] - 'a'];
            if (pass[cur] == 1) {
                out.push_back(w.substr(0, k + 1));
                break;
            }
        }
    }
    return out;
}
