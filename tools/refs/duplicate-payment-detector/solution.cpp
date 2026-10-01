#include <map>
#include <tuple>
#include <vector>
using namespace std;

vector<int> duplicates(vector<vector<int>>& payments, int window) {
    map<tuple<int, int, int>, int> last;
    vector<int> out;
    for (int i = 0; i < (int)payments.size(); i++) {
        auto& p = payments[i];
        auto key = make_tuple(p[1], p[2], p[3]);
        auto it = last.find(key);
        if (it != last.end() && p[0] - it->second <= window) out.push_back(i);
        last[key] = p[0];
    }
    return out;
}
