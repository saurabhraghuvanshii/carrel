#include <string>
#include <vector>
using namespace std;

vector<int> nearMatches(vector<string>& words, vector<string>& queries) {
    vector<int> out;
    for (auto& q : queries) {
        int count = 0;
        for (auto& w : words) {
            if (w.size() != q.size()) continue;
            int diff = 0;
            for (size_t i = 0; i < w.size(); i++) diff += w[i] != q[i];
            if (diff <= 1) count++;
        }
        out.push_back(count);
    }
    return out;
}
