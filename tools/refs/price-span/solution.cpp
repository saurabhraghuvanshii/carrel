#include <vector>
using namespace std;

vector<int> spans(vector<int>& prices) {
    vector<int> span(prices.size()), higher;
    for (int i = 0; i < (int)prices.size(); i++) {
        while (!higher.empty() && prices[higher.back()] <= prices[i]) higher.pop_back();
        span[i] = higher.empty() ? i + 1 : i - higher.back();
        higher.push_back(i);
    }
    return span;
}
